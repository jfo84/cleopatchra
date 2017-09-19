require_relative './data_mixin'

class Pull
  include Virtus.model
  include DataMixin

  def record
    connection.exec("INSERT INTO pulls (id, data)
                    VALUES (#{id}, #{data})")
  end

  def url
    data_hash['url']
  end
end